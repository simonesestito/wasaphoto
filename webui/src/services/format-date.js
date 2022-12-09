/**
 * Format a given date and time
 * @param {Date|string} date Datetime to format
 * @returns {string} Formatted datetime
 */
export function formatDate(date) {
	if (typeof date === 'string') {
		date = new Date(date);
	}

	return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}
