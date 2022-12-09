import { api } from './axios';
import { getCurrentUID } from './auth-store';
import {handleApiError, NotFoundError} from "./api-errors";

export const LikesService = Object.freeze({
    /**
     * Like a photo
     * @param {string} photoId Photo to like
     */
    async likePhoto(photoId) {
        const response = await api.put(`/photos/${photoId}/likes/${getCurrentUID()}`);

        switch (response.status) {
            case 200: case 201: return;
            case 404: throw new NotFoundError('Photo to like not found');
            default: handleApiError(response);
        }
    },

    /**
     * Unlike a photo
     * @param {string} photoId Photo to unlike
     */
     async unlikePhoto(photoId) {
        const response = await api.delete(`/photos/${photoId}/likes/${getCurrentUID()}`);

        switch (response.status) {
            case 204: return;
            case 404: throw new NotFoundError('Photo to unlike not found');
            default: handleApiError(response);
        }
    }
});
