import {ConflictError, handleApiError, NotFoundError} from "./api-errors";
import api from "./axios";
import {getCurrentUID} from "./auth-store";

export const FollowService = Object.freeze({
    /**
     * List someone's followers
     * @param {string} userId
     * @param {string?} pageCursor
     */
    async listFollowers(userId, pageCursor) {
        let apiPath = `/users/${userId}/followers/`;
        if (pageCursor) {
            apiPath += '?pageCursor=' + encodeURIComponent(pageCursor);
        }

        const response = await api.get(apiPath);

        switch (response.status) {
            case 200: return response.data;
            case 404: throw new NotFoundError('User not found');
            default: handleApiError(response);
        }
    },

    /**
     * List someone's followings
     * @param {string} userId
     * @param {string?} pageCursor
     */
     async listFollowings(userId, pageCursor) {
        let apiPath = `/users/${userId}/followings/`;
        if (pageCursor) {
            apiPath += '?pageCursor=' + encodeURIComponent(pageCursor);
        }

        const response = await api.get(apiPath);

        switch (response.status) {
            case 200: return response.data;
            case 404: throw new NotFoundError('User not found');
            default: handleApiError(response);
        }
    },

    /**
     * Unfollow a user
     * @param {string} followedId
     */
     async unfollowUser(followedId) {
        const response = await api.delete(`/users/${getCurrentUID()}/followings/${followedId}`);

        switch (response.status) {
            case 204: return;
            default: handleApiError(response);
        }
    },

    /**
     * Follow a user
     * @param {string} followedId
     */
     async followUser(followedId) {
        const response = await api.put(`/users/${getCurrentUID()}/followings/${followedId}`);

        switch (response.status) {
            case 200: case 201: return;
            case 404: throw new NotFoundError('User to follow not found');
			case 409: throw new ConflictError(response.data);
            default: handleApiError(response);
        }
    }
});
