const AUTH_TOKEN_KEY = 'auth-token';

/**
 * Save the new auth token received
 * @param {string} authToken
 */
export function saveAuthToken(authToken) {
    localStorage.setItem(AUTH_TOKEN_KEY, authToken);
}

/**
 * Retrieve the saved auth token, if any.
 * @returns {string|null} Saved auth token, if any
 */
export function getAuthToken() {
    return localStorage.getItem(AUTH_TOKEN_KEY);
}

/**
 * Retrieve the ID of the current user.
 * It may be different from the auth token, not in this case.
 *
 * @returns {string|null} Current logged-in user ID
 */
export function getCurrentUID() {
    return getAuthToken();
}
