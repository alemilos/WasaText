import axios from "axios";

const axiosInstance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});

const abortController = new AbortController();

axiosInstance.interceptors.request.use(
	(config) => {
		config.signal = abortController.signal;
		config.headers["Authorization"] = config.identifier; // just use the user id
		return config;
	},
	(error) => {
		console.log("interceptor found error: ", error);
		abortController.abort();
		Promise.reject(error);
	}
);

axiosInstance.interceptors.response.use();

export default axiosInstance;
