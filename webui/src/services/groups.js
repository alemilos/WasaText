import { Endpoints } from "../config/endpoints";
import { del, put, post } from "./axios";

export function addToGroup() {
	const { groupId } = boh;
	const url = Endpoints.addToGroup.replace(":id", groupId);
	return post(url);
}

export function leaveGroup() {
	const { groupId } = boh;
	const url = Endpoints.leaveGroup.replace(":id", groupId);
	return del(url);
}

export function setGroupName() {
	const { groupId } = boh;
	const url = Endpoints.setGroupName.replace(":id", groupId);
	return put(url);
}

export function setGroupPhoto() {
	const { groupId } = boh;
	const url = Endpoints.setGroupPhoto.replace(":id", groupId);
	return put(url);
}

export function createGroup(name, members) {
	return post(Endpoints.createGroup, { name, members });
}
