import { endpoints } from "../config/endpoints";
import { post } from "./axios";

export function saveUserId(id) {
	localStorage.setItem("user-id", id);
}

export function getUserId() {
	return localStorage.getItem("user-id");
}

export function clearUserId() {
	localStorage.removeItem("user-id");
}

export function doLogin(username) {
	console.log(username);
	username = username.trim();
	if (!username || username === "") {
		throw new Error("Invalid Username");
	}
	return post(endpoints.doLogin, { username });
}
