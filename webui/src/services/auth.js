import { handleApiError } from './api-errors';
import { saveAuthToken } from './auth-store';
import api from './axios';

/**
 * @typedef {Object} UserLoginResult
 * User Login Result
 * @property {string} userId Auth token and user ID
 * @property {boolean} isNewUser Indicates whether the operation was a login or a sign up
 */

export const AuthService = Object.freeze({
    /**
     * Login or sign up
     * @param {string} username User username
	 * @param {boolean} keepSignedIn Keep me signed in
     * @returns {Promise<UserLoginResult>}
     */
    async doLogin(username, keepSignedIn) {
        const response = await api.post('/session', {
            username: username,
        });

        // Save new auth token
        if (response.status >= 200 && response.status < 300 && response.data.userId) {
            saveAuthToken(response.data.userId, keepSignedIn);
        }

        if (response.status === 200) {
            return {
                userId: response.data.userId,
                isNewUser: false,
            };
        } else if (response.status === 201) {
            return {
                userId: response.data.userId,
                isNewUser: true,
            };
        } else {
            handleApiError(response);
        }
    },

    /**
     * Logout current user
     */
    logout() {
        saveAuthToken(null);
    }
});
