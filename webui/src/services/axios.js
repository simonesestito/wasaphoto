import axios from "axios";

export const api = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

export default api;

