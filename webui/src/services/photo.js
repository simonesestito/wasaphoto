import api from "./axios";
import {BadRequestError, handleApiError, NotFoundError, ThirdPartyError} from "./api-errors";

export const PhotosService = Object.freeze({
	/**
	 * List someone's photos
	 * @param {string} userId User to see photos of
	 * @param {string|null} pageCursor Current page cursor
	 */
    async listUserPhotos(userId, pageCursor) {
		let apiPath = `/users/${userId}/photos/`;
		if (pageCursor) {
			apiPath += '?pageCursor=' + pageCursor;
		}

		const response = await api.get(apiPath);

		switch (response.status) {
			case 200: return response.data;
			case 404: throw new NotFoundError('User not found');
			default: handleApiError(response);
		}
	},

	/**
	 * Upload photo
	 * @param {File} photoFile
	 * @param {Function} onProgress Progress callback
	 */
	async uploadPhoto(photoFile, onProgress) {
		const response = await api.post('/photos/', await photoFile.arrayBuffer(), {
			timeout: 60000, // Enlarge timeout for photo upload, check if enough on server side
			onUploadProgress: progressEvent => onProgress(100.0 * progressEvent.loaded / progressEvent.total),
		});

		switch (response.status) {
			case 201: return response.data;
			case 415: throw new BadRequestError('Selected photo file cannot be processed as an image');
			case 503: throw new ThirdPartyError();
			default: handleApiError(response);
		}
	},

	/**
	 * Delete a photo
	 * @param {string} photoId Photo to delete
	 */
	async deletePhoto(photoId) {
		const response = await api.delete(`/photos/${photoId}`);

		if (response.status !== 204) {
			handleApiError(response);
		}
	}
});
