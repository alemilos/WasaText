import { Endpoints } from "../config/endpoints";
import { get, post, put } from "./axios";

export function setMyUsername() {
	return put(Endpoints.setMyUsername);
}

export function setMyPhoto() {
	return post(Endpoints.setMyPhoto);
}

export function getUsers() {
	return get(Endpoints.getUsers);
}
