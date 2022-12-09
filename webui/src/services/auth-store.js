const AUTH_TOKEN_KEY = 'auth-token';

/**
 * Save the new auth token received
 * @param {string|null} authToken
 * @param {boolean} [keepSignedIn] Keep me signed in
 */
export function saveAuthToken(authToken, keepSignedIn) {
	const storage = keepSignedIn ? localStorage : sessionStorage;
	const resetStorage = keepSignedIn ? sessionStorage : localStorage;

    if (authToken) {
		storage.setItem(AUTH_TOKEN_KEY, authToken);
	} else {
		storage.removeItem(AUTH_TOKEN_KEY);
	}

	resetStorage.removeItem(AUTH_TOKEN_KEY);
}

/**
 * Retrieve the saved auth token, if any.
 * @returns {string|null} Saved auth token, if any
 */
export function getAuthToken() {
	const persistingToken = localStorage.getItem(AUTH_TOKEN_KEY);
	if (persistingToken) {
		return persistingToken;
	}

	return sessionStorage.getItem(AUTH_TOKEN_KEY);
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
