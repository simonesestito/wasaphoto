/**
 * Handle unexpected errors that may arrive as a response to an API call
 * @param {import("axios").AxiosResponse<any, any>} response The response with an unexpected result
 */
export function handleApiError(response) {
    console.error('Unexpected error on API call');
    console.error(response.statusText);
    console.error(response.data);

    switch (response.status) {
        case 400: throw new BadRequestError(null);
        case 401: throw new AuthError();
        case 403: throw new ForbiddenError();
		case 404: throw new NotFoundError('Item not found');
		case 409: throw new ConflictError(null);
		case 413: throw new TooLargeError();
        case 500: throw new ServerError();
		case 503: throw new ThirdPartyError();
		default: throw new Error('Unexpected error received from server: ' + response.status);
    }
}

export class BadRequestError extends Error {
	/**
	 * Bad request, with an optional custom message
	 * @param {string|null} message Custom error message
	 */
    constructor(message) {
        super(message || "Unexpected error on client side");
    }
}

export class ServerError extends Error {
    constructor() {
        super("Unexpected error on server side");
    }
}

export class ConflictError extends Error {
	/**
	 * Conflict error, with an optional custom message
	 * @param {string|null} message Custom error message
	 */
	constructor(message) {
		super(message || 'A conflict occurred');
	}
}

export class AuthError extends Error {
    constructor() {
        super("No auth or invalid auth was provided: try to login again");
    }
}

export class ForbiddenError extends Error {
    constructor() {
        super("You are not authorized to perform this action");
    }
}

export class NotFoundError extends Error {
    constructor(message) {
        super(message);
    }
}

export class ThirdPartyError extends Error {
	constructor() {
		super("This operation is unavailable at the moment because a third-party service is unavailable");
	}
}

export class TooLargeError extends Error {
	constructor() {
		super("The file you're uploading is too big");
	}
}
