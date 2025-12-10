import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { target: 50, duration: '30s' },
        { target: 50, duration: '1m' },
        { target: 0, duration: '30s' },
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<300'],
    },
};

let minTaskId = 1;
let maxTaskId = 42397;

const randomTaskId = () => {
    return Math.floor(Math.random() * (maxTaskId - minTaskId + 1)) + minTaskId;
};

export default function () {
    const res = http.del(`http://localhost:1919/api/todo/tasks/${randomTaskId()}`);

    check(res, {
        'status is 204': (r) => r.status === 204
    });
}
