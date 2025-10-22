import { inject } from "vue";

export function useModal() {
	const openModal = inject("openModal");
	const closeModal = inject("closeModal");
	const modal = inject("modal");

	if (!openModal || !closeModal || !modal) {
		throw new Error("useModal must be used within ModalProvider");
	}

	return { openModal, closeModal, modal };
}
