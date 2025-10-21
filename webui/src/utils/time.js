const MONTHS_IT = [
	"gen",
	"feb",
	"mar",
	"apr",
	"mag",
	"giu",
	"lug",
	"ago",
	"set",
	"ott",
	"nov",
	"dic",
];

/**
 * Format a UTC ISO date string into human-friendly Italian format
 * @param isoDate string like "2025-10-21T11:22:59Z"
 */
export function formatDate(isoDate) {
	if (!isoDate) return "";

	const date = new Date(isoDate);
	const now = new Date();

	const diffMs = now.getTime() - date.getTime();
	const diffSeconds = Math.floor(diffMs / 1000);

	// less than a minute ago
	if (diffSeconds < 60) return "ora";

	// today
	const isToday = now.toDateString() === date.toDateString();
	if (isToday) {
		return date.toLocaleTimeString("it-IT", {
			hour: "2-digit",
			minute: "2-digit",
		});
	}

	// yesterday
	const yesterday = new Date(now);
	yesterday.setDate(now.getDate() - 1);
	if (yesterday.toDateString() === date.toDateString()) {
		return "ieri";
	}

	// other days
	const day = date.getDate();
	const month = MONTHS_IT[date.getMonth()];
	const year = date.getFullYear();

	return `${day} ${month} ${year}`;
}
