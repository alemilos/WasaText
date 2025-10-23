import { Endpoints } from "../config/endpoints";
import { get, post, postForm, put } from "./axios";

export function setMyUsername(username) {
	return put(Endpoints.setMyUsername, { username });
}

export function setMyPhoto(imageFile) {
	const formData = new FormData();
	formData.append("file", imageFile);

	return postForm(Endpoints.setMyPhoto, formData);
}

export function getUsers() {
	return get(Endpoints.getUsers);
}
