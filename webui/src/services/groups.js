import { Endpoints } from "../config/endpoints";
import { del, put, post, putForm } from "./axios";

export function addToGroup(groupId, userId) {
	const url = Endpoints.addToGroup.replace(":id", groupId);
	return post(url, { userId });
}

export function leaveGroup(groupId, userId) {
	const url = Endpoints.leaveGroup.replace(":id", groupId);
	return del(url, { userId });
}

export function setGroupName(groupId, name) {
	const url = Endpoints.setGroupName.replace(":id", groupId);
	return put(url, { name });
}

export function setGroupPhoto(groupId, imageFile) {
	const formData = new FormData();
	formData.append("file", imageFile);

	const url = Endpoints.setGroupPhoto.replace(":id", groupId);
	return putForm(url, formData);
}

export function createGroup(name, members) {
	return post(Endpoints.createGroup, { name, members });
}
