// selectTab shows tab index of a tab group inside root.
export function selectTab(group, index, root) {
	(root || document).querySelectorAll(`[data-group="${group}"]`).forEach((el) => {
		el.classList.toggle("active", el.dataset.tab === String(index));
	});
}
