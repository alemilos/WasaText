export function getError(res, fallback = "Errore") {
	return (
		res?.error ||
		res.response?.data?.message ||
		res.response?.data ||
		fallback
	);
}
