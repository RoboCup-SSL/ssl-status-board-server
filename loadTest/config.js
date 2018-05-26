import ws from "k6/ws";
import {check, sleep} from "k6";

let minConnectionDuration = 10;
let maxConnectionDuration = 60;
let users = 400;
let iterationsPerUser = 2;
let maxPauseTime = 30;

// noinspection JSUnusedGlobalSymbols
export let options = {
    noConnectionReuse: true,
    iterations: users * iterationsPerUser,
    vus: users,
};

// noinspection JSUnusedGlobalSymbols
export default function () {
    const url = "wss://tigers-mannheim.de/ssl-vision/field-a/subscribe";
    const params = {};

    let packagesReceived = 0;

    sleep(Math.random() * maxPauseTime);

    const connectionDuration = Math.round(minConnectionDuration + Math.random() * (maxConnectionDuration - minConnectionDuration)) * 1000;

    const res = ws.connect(url, params, function (socket) {

        socket.on('message', function (data) {
            packagesReceived++;
        });

        socket.setTimeout(function () {
            socket.close();
        }, connectionDuration);
    });

    check(res, {"status is 101": (r) => r && r.status === 101});
    check(packagesReceived, {"received packages": (r) => r > 0});
}
