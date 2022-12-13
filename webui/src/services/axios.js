import axios from "axios";
import {getCurrentUID, saveAuthToken} from "./auth-store";
import router from "../router";

export const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
	validateStatus: (_) => true,
});

api.interceptors.request.use(config => {
	config.headers['Authorization'] = 'Bearer ' + getCurrentUID();
	return config;
});

api.interceptors.response.use(async response => {
	if (response && response.status === 401) {
		// Logout!
		saveAuthToken(null);
		await router.redirectToLogin();
	}
	return response;
});

export default api;

