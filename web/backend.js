const sendRequest = (path, options) => {
    return fetch(path, options)
        .then(async res => {
            return res.ok ? res.json() : Promise.reject(await res.json());
        })
        .catch(err => {
            // console.error(err);
            throw err;
        });
};

const apiRequest = async (method, path, data) => {
    const options = {
        method: method,
    };
    if (options.method !== 'GET') {
        options.headers = {
            'Content-Type': 'application/json',
        };
        options.body = JSON.stringify(data);
    }
    return sendRequest(path, options);
};

const apiUploadRequest = async (path, file) => {
    const options = {
        method: 'POST',
    };

    const form = new FormData();
    form.append('file', new Blob([file]));

    options.body = form;

    return sendRequest(path, options);
};

export default class Backend {
    baseUrl = 'http://localhost:8000'

    constructor(baseUrl) {
        this.baseUrl = baseUrl;
    }

    uploadSchema(file) {
        return apiUploadRequest(this.baseUrl + '/api/schemas', file);
    }

    fetchSchemas() {
        return apiRequest('GET', this.baseUrl + '/api/schemas');
    }

    publishMessage(data) {
        return apiRequest('POST', this.baseUrl + '/api/message', data);
    }

    connectMessagesStream(next) {
        const client = new EventSource(this.baseUrl + '/api/stream');
        client.addEventListener('message', (event) => {
            next(JSON.parse(event.data));
        });
    }

    checkVersion() {
        return apiRequest('GET', this.baseUrl + '/api/check_version');
    }
}
