import { endpoints } from "../config/endpoints";
import { del, put, post } from "./axios";

export function addToGroup() {
	const { groupId } = boh;
	const url = endpoints.addToGroup.replace(":id", groupId);
	return post(url);
}

export function leaveGroup() {
	const { groupId } = boh;
	const url = endpoints.leaveGroup.replace(":id", groupId);
	return del(url);
}

export function setGroupName() {
	const { groupId } = boh;
	const url = endpoints.setGroupName.replace(":id", groupId);
	return put(url);
}

export function setGroupPhoto() {
	const { groupId } = boh;
	const url = endpoints.setGroupPhoto.replace(":id", groupId);
	return put(url);
}

export function createGroup() {
	return post(endpoints.createGroup);
}
