import axios from "axios";
import router from "../router";
import { getUserId } from "./auth";

const axiosInstance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});

axiosInstance.interceptors.request.use((config) => {
	const abortController = new AbortController();
	config.signal = abortController.signal;
	config.headers["Authorization"] = "Bearer " + getUserId(); // just use the user id
	return config;
});

axiosInstance.interceptors.response.use(
	async (response) => {
		if (response && response.status === 401) {
			// logout the user when NOT AUTHORIZED
			await router.logoutAndRedirect();
		}

		return response;
	},
	async (error) => {
		console.log(error);
		const status = error?.response?.status;
		if (status === 401) {
			// logout the user when NOT AUTHORIZED
			await router.logoutAndRedirect();
		}

		return error;
	}
);

export async function get(url, config = {}) {
	return axiosInstance.get(url, { ...config }).then((response) => response);
}

export async function post(url, data = {}, config = {}) {
	return axiosInstance.post(url, data, { ...config });
}

export async function put(url, data = {}, config = {}) {
	if (typeof data === "string") data = JSON.parse(data);

	return axiosInstance
		.put(url, { ...data }, { ...config })
		.then((response) => response);
}

export async function postForm(url, form, config = {}) {
	return axiosInstance
		.postForm(url, form, {
			transformRequest: (formData) => formData,
			...config,
		})
		.then((response) => response);
}

export async function del(url, config = {}) {
	return axiosInstance
		.delete(url, { ...config })
		.then((response) => response);
}

export default axiosInstance;
