import ws from "k6/ws";
import http from "k6/http";
import {check, sleep} from "k6";

let minConnectionDuration = 30;
let maxConnectionDuration = 120;
let users = 500;
let iterationsPerUser = 3;
let maxPauseTime = 60;

// noinspection JSUnusedGlobalSymbols
export let options = {
    noConnectionReuse: true,
    iterations: users * iterationsPerUser,
    vus: users,
};

// noinspection JSUnusedGlobalSymbols
export default function () {
    const homeUrl = "https://tigers-mannheim.de/status-board/";
    const visionUrl = "wss://tigers-mannheim.de/ssl-vision/field-a/subscribe";
    const params = {};

    const connectionDuration = Math.round(minConnectionDuration + Math.random() * (maxConnectionDuration - minConnectionDuration)) * 1000;

    let visionPackagesReceived = 0;

    sleep(Math.random() * maxPauseTime);

    let homeRes = http.get(homeUrl);
    check(homeRes, {"entry status is 200": (r) => r && r.status === 200});

    sleep(Math.random() * 2);

    const visionRes = ws.connect(visionUrl, params, function (socket) {

        socket.on('message', function (data) {
            visionPackagesReceived++;
        });

        socket.setTimeout(function () {
            socket.close();
        }, connectionDuration);
    });

    check(visionRes, {"status is 101": (r) => r && r.status === 101});
    check(visionPackagesReceived, {"received vision packages": (r) => r > 0});
}
