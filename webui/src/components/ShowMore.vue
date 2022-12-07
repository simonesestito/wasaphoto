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
	<div class="row justify-content-center">
		<button type="button" class="btn btn-outline-light" @click="loadMore">Load more</button>
	</div>
</template>
