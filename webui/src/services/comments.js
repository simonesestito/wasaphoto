import { handleApiError, NotFoundError } from './api-errors';
import { api } from './axios';

export const CommentsService = Object.freeze({
    /**
     * Get paginated photo comments
     * @param {string} photoId
     * @param {string?} pageCursor
     */
    async getPhotoComments(photoId, pageCursor) {
        let apiPath = `/photos/${photoId}/comments/`;
        if (pageCursor) {
            apiPath += '?pageCursor=' + encodeURIComponent(pageCursor);
        }

        const response = await api.get(apiPath);

        switch (response.status) {
            case 200: return response.data;
            case 404: throw new NotFoundError('Required post does not exist');
            default: handleApiError(response);
        }
    },

    /**
     * Comment a photo
     * @param {string} photoId Photo where to publish a comment to
     * @param {string} text Comment text
     */
     async commentPhoto(photoId, text) {
        const response = await api.post(`/photos/${photoId}/comments/`, {
            text: text,
        });

        switch (response.status) {
            case 201: return response.data;
            case 404: throw new NotFoundError('Required post does not exist');
            default: handleApiError(response);
        }
    },

    /**
     * Delete a comment from a photo
     * @param {string} photoId Photo with the comment to remove
     * @param {string} commentId Comment ID to delete
     */
     async uncommentPhoto(photoId, commentId) {
        const response = await api.delete(`/photos/${photoId}/comments/${commentId}`);

        switch (response.status) {
            case 204: return;
            case 404: throw new NotFoundError('Required post does not exist');
            default: handleApiError(response);
        }
    }
});
