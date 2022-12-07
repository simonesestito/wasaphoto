import api from "./axios";
import {ConflictError, handleApiError, NotFoundError} from "./api-errors";
import {getCurrentUID} from "./auth-store";

export const UsersService = Object.freeze({
	/**
	 * Search users by username, performing a text search
	 * @param {string} username Username or partial username to search
	 * @param {string|null} pageCursor Search page cursor
	 * @returns Promise<any>
	 */
	async searchUsers(username, pageCursor) {
		return _searchUsers(username, {exactMatch: false}, pageCursor);
	},

	/**
	 * Get a user by a given username, if exists
	 * @param username
	 * @returns {Promise<Object | null>} User with this username, if any
	 */
	async getByUsername(username) {
		const pageResult = await _searchUsers(username, {exactMatch: true}, null);
		const foundUser = pageResult.pageData.find(user => user.username === username);
		if (!foundUser) {
			throw new NotFoundError('User not found');
		}
		return foundUser;
	},

	/**
	 * Get a user profile by ID
	 * @param userId ID of the user to fetch
	 */
	async getUserProfile(userId) {
		const response = await api.get(`/users/${userId}`);
		if (response.status === 200) {
			return response.data;
		} else {
			handleApiError(response);
		}
	},

	/**
	 * Set my user details.
	 * @param {string} userDetails.name
	 * @param {string} userDetails.surname
	 * @param {string} userDetails.username
	 * @returns {Promise<void>}
	 */
	async setMyDetails(userDetails) {
		const response = await api.put(`/users/${getCurrentUID()}`, userDetails);

		switch (response.status) {
			case 200: return response.data;
			case 409: throw new ConflictError("Username already taken");
			default: handleApiError(response);
		}
	},

	/**
	 * Set my username.
	 * @param {string} username
	 * @returns {Promise<void>}
	 */
	async setMyUsername(username) {
		const response = await api.put(`/users/${getCurrentUID()}/username`, username);

		switch (response.status) {
			case 200: return response.data;
			case 409: throw new ConflictError("Username already taken");
			default: handleApiError(response);
		}
	}
});

/**
 * Private function.
 * Search users by username, either with exact username or other ways.
 * @param {string} username username to search, or partial, according to 'exactMatch'
 * @param {boolean} exactMatch Indicates the type of search to perform
 * @param {string|null} pageCursor Page cursor of the current search, if any
 * @returns Promise<any>
 */
async function _searchUsers(username, {exactMatch}, pageCursor) {
	const apiPath = Object.entries({
		username, exactMatch, pageCursor
	}).filter(([_, value]) => value !== null)
		.map(([key, value]) => [key, value == null ? null : value.toString()])
		.map(entry => entry.map(encodeURIComponent))
		.reduce((acc, [key, value]) => `${acc}&${key}=${value}`, '/users/?');

	const response = await api.get(apiPath);

	if (response.status === 200) {
		return response.data;
	} else {
		handleApiError(response);
	}
}
