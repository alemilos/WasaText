import { Endpoints } from "../config/endpoints";
import { LocalStorage } from "../config/localStorage";
import { post } from "./axios";

export function saveUserId(id) {
	localStorage.setItem(LocalStorage.userId, id);
}

export function getUserId() {
	return localStorage.getItem(LocalStorage.userId);
}

export function clearUserId() {
	localStorage.removeItem(LocalStorage.userId);
}

export function doLogin(username) {
	username = username.trim();
	if (!username || username === "") {
		throw new Error("Invalid Username");
	}

	return post(Endpoints.doLogin, { username });
}
