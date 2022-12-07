<script>
export default {
	name: "ShowMore",
	emits: ['loadMore'],
	methods: {
		loadMore() {
			this.$emit('loadMore');
		},
	},
};
</script>

<template>
	<div class="row-cols-1 mt-4">
		<button type="button" class="btn btn-outline-secondary" @click="loadMore">Load more</button>
	</div>
</template>
