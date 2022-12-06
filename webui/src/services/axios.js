import axios from "axios";
import {getCurrentUID} from "./auth-store";

export const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
	validateStatus: (_) => true,
});

api.interceptors.request.use(function (config) {
	config.headers.Authorization = 'Bearer ' + getCurrentUID();
	return config;
});

export default api;

