import { reactive, computed } from "vue";

const state = reactive({
	replyingTo: null,
});

export const replyStore = {
	setReply(message) {
		state.replyingTo = message;
	},

	clearReply() {
		state.replyingTo = null;
	},

	replyingTo: computed(() => state.replyingTo),
};
