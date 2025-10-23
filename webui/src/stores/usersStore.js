import { reactive } from "vue";
import { getUsers } from "../services/users";

const users = reactive({}); // key: userId, value: user object

async function loadUsers() {
	try {
		const res = await getUsers();
		if (res.status === 200) {
			// clear old data
			for (const key in users) delete users[key];

			// add new users
			res.data.users.forEach((user) => {
				users[user.userId] = user;
			});
		}
	} catch (err) {
		console.error("Failed to load users:", err);
	}
}

function getUser(userId) {
	return users[userId];
}

function getUsername(userId) {
	return users[userId]?.username ?? "Unknown User";
}

function getPhotoUrl(userId) {
	const path = users[userId]?.photoPath ?? null;
	if (!path) return null;
	// Add a version query string to bypass caching
	return `${path}?v=${Date.now()}`;
}

function updatePhotoPath(userId, photoPath) {
	if (!users[userId]) {
		return;
	}
	// Add a version query string to bypass caching
	users[userId].photoPath = `${photoPath}?v=${Date.now()}`;
}

function updateUsername(userId, username) {
	if (!users[userId]) {
		return;
	}
	users[userId].username = username;
}

export const usersStore = {
	users,
	loadUsers,
	getUser,
	getUsername,
	getPhotoUrl,
	updatePhotoPath,
	updateUsername,
};
