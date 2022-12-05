import api from "./axios";
import {getCurrentUID} from "./auth-store";
import {handleApiError} from "./api-errors";


export const StreamService = Object.freeze({
    /**
     * Get my own stream, paginated
     * @param {string|null} pageCursor
     */
    async getMyStream(pageCursor) {
        const response = await api.get(`/users/${getCurrentUID()}/stream`);

		if (response.status === 200) {
			return response.data;
		} else {
			handleApiError(response);
		}
    }
});
