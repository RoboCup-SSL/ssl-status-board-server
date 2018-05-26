import ws from "k6/ws";
import {check} from "k6";

let connectionDuration = 60;

// noinspection JSUnusedGlobalSymbols
export let options = {
    noConnectionReuse: true,
};

// noinspection JSUnusedGlobalSymbols
export default function () {
    const url = "wss://tigers-mannheim.de/ssl-vision/field-a/subscribe";
    const params = {};

    let packagesReceived = 0;

    const res = ws.connect(url, params, function (socket) {
        socket.on('open', function () {
            //console.log('connected');
        });

        socket.on('message', function (data) {
            packagesReceived++;
            //console.log("Message received: ", data);
        });

        socket.on('close', function () {
            //console.log('disconnected');
        });

        socket.setTimeout(function () {
            socket.close();
        }, connectionDuration * 1000);
    });

    check(res, {"status is 101": (r) => r && r.status === 101});
    check(packagesReceived, {"received packages": (r) => r > 0});
}
