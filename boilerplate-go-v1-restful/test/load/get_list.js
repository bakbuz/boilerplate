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

export default function () {
    const res = http.get(`http://localhost:1919/api/todo/lists/${randomListId()}`);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response body': (r) => {
            const result = JSON.parse(r.body);
            return result.list && result.list.id > 0;
        },
    });
}
