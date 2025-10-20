import { endpoints } from "../config/endpoints";
import { get, post, put } from "./axios";

export function setMyUsername() {
	return put(endpoints.setMyUsername);
}

export function setMyPhoto() {
	return post(endpoints.setMyPhoto);
}

export function getUsers() {
	return get(endpoints.getUsers);
}
