import api from "./axios";
import {getCurrentUID} from "./auth-store";
import {ConflictError, handleApiError, NotFoundError} from "./api-errors";

export const BanService = Object.freeze({
	/**
	 * Ban a user
	 * @param blockedId ID of the user to ban
	 */
	async banUser(blockedId) {
		const response = await api.put(`/users/${getCurrentUID()}/bannedPeople/${blockedId}`);

		switch (response.status) {
			case 200: case 201: return;
			case 404: throw new NotFoundError("User to ban not found");
			case 409: throw new ConflictError("You cannot ban yourself");
			default: handleApiError(response);
		}
	},

	/**
	 * Unban a user
	 * @param blockedId ID of the user to ban
	 */
	async unbanUser(blockedId) {
		const response = await api.delete(`/users/${getCurrentUID()}/bannedPeople/${blockedId}`);

		if (response.status !== 204) {
			handleApiError(response);
		}
	},
});
