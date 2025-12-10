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

let minListId = 1;
let maxListId = 51139;

const randomListId = () => {
    return Math.floor(Math.random() * (maxListId - minListId + 1)) + minListId;
};

const randomTaskNum = () => {
    return Math.floor(Math.random() * 1000000);
};

export default function () {
    const body = {
        name: `load_test${randomTaskNum()}`,
    };

    const res = http.put(`http://localhost:1919/api/todo/lists/${randomListId()}/tasks`, JSON.stringify(body), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        'status is 201': (r) => r.status === 201,
        'response body': (r) => {
            const result = JSON.parse(r.body);
            return result.id > 0;
        },
    });
}
