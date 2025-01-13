import axios from "axios";
import router from "../router";
import { getUserId } from "./auth";

const axiosInstance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});


axiosInstance.interceptors.request.use(
	(config) => {
		const abortController = new AbortController();
		config.signal = abortController.signal;
		config.headers["Authorization"] = "Bearer " + getUserId(); // just use the user id
		return config;
	},
);

axiosInstance.interceptors.response.use(async response => {
	if (response && response.status === 401) {
		// logout the user when NOT AUTHORIZED
		await router.logoutAndRedirect();
	}
	return response;
});

export default axiosInstance;
