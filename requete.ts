//const BASE_URL = 'http://localhost:9124' // webpack
//const BASE_URL = 'http://localhost:3333'  // chi
const BASE_URL = 'http://localhost:5173'    // vite

export class ValidationError extends Error {
    public details: any;

    constructor(message: string, details?: any) {
        super(message);
        this.details = details;
        // Pour une meilleure compatibilité avec instanceof
        Object.setPrototypeOf(this, ValidationError.prototype);
    }
}

export class UnprocessableError extends Error {
    public details: any;

    constructor(message: string, details?: any) {
        super(message);
        this.details = details;
        // Pour une meilleure compatibilité avec instanceof
        Object.setPrototypeOf(this, UnprocessableError.prototype);
    }
}

/*
Evite d'avoir 2 requetes qui tentent d'acquerir un nouveau token d'accés en meme temps.
*/
let isRefreshing = false;

function waitForRefreshing(timeout = 10000) {
    return new Promise((resolve, reject) => {
        const startTime = Date.now();
        const interval = setInterval(() => {
            if (!isRefreshing) {
                clearInterval(interval);
                resolve(null);
            } else {
                console.log("attend")
                if (Date.now() - startTime >= timeout) {
                    clearInterval(interval);
                    reject(new Error("Timeout waiting for isRefreshing to be false"));
                }
            }
        }, 100); // Vérifier toutes les 100 ms
    });
}

export async function api<T>(path: string, method?: string, body?: any): Promise<T> {

    const init: RequestInit = {
        body, method
    }


    var response = await fetch(`${BASE_URL}${path}`, init);

    if (response.status === 401) {
        await waitForRefreshing()

        isRefreshing = true
        const respRefresh = await fetch(`${BASE_URL}/api/v0/auth/refresh`)
        isRefreshing = false

        if (respRefresh.status === 204) {
            response = await fetch(`${BASE_URL}${path}`, init);
        }
    }

    if (!response.ok) {
        const contentType = response.headers.get('Content-Type');
        let errorBody: any = undefined;

        if (contentType && contentType.includes('application/json')) {
            try {
                errorBody = await response.json();
            } catch {
                errorBody = undefined;
            }
        } else {
            // Essaye de lire le texte brut
            const text = await response.text();
            if (text) errorBody = text;
        }

        if (response.status === 400) {
            throw new ValidationError(response.statusText, errorBody);
        }

        if (response.status === 422) {
            throw new UnprocessableError(response.statusText, errorBody);
        }

        throw new Error(response.statusText + (errorBody ? `: ${errorBody}` : ""));
    }

    if (response.status === 204) {
        return new Promise<any>((resolve) => {
            resolve(null);
        });
    }

    const contentTypeHeader = response.headers.get('Content-Type');
    if (contentTypeHeader == null) {
        throw new Error("pas de type dans la reponse");
    } else if (contentTypeHeader.includes('application/json')) {
        return await response.json() as T;
    } else {
        return response as T
    }

}

