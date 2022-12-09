<template>
	<div class="modal fade" tabindex="-1">
		<div class="modal-dialog">
			<form class="modal-content" @submit.prevent="publishComment">
				<div class="modal-header">
					<h5 class="modal-title">Post a comment</h5>
					<button ref="close" type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<input type="text" minlength="1" maxlength="256" required v-model="commentText">
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
					<input class="btn btn-primary" type="submit" value="Publish Comment">
				</div>
			</form>
		</div>
	</div>
</template>

<script>
export default {
	name: "NewCommentModal",
	emits: ['comment'],
	data() {
		return {
			commentText: '',
		};
	},
	methods: {
		publishComment() {
			this.$emit('comment', this.commentText);
			this.$refs.close.click();
		},
	},
}
</script>
