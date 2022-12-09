<template>
	<div class="mt-5">
		<UserNameHeader :user="comment.author"/>
		<div>
			<p>{{ comment.text }}</p>
			<span><i>{{ formatDate(comment.publishDate) }}</i></span>
		</div>
		<button v-if="comment.author.id === getCurrentUID()"
				@click="() => deleteComment(comment.id)"
				type="button"
				:disabled="loading"
				class="btn btn-outline-danger">
			Delete
		</button>
	</div>
</template>

<script>
import {CommentsService} from "../services";
import {formatDate} from "../services/format-date";
import {getCurrentUID} from "../services/auth-store";
import UserNameHeader from "./UserNameHeader.vue";

export default {
	name: "CommentListItem",
	components: {UserNameHeader},
	props: ['comment'],
	emits: ['error', 'refresh'],
	data() {
		return {
			loading: false,
		};
	},
	methods: {
		async deleteComment(commentId) {
			this.loading = true;
			this.$emit('error', '');

			try {
				await CommentsService.uncommentPhoto(this.$route.params.photoId, commentId);
				this.$emit('refresh');
			} catch (err) {
				this.$emit('error', err.toString());
			} finally {
				this.loading = false;
			}
		},
		formatDate(date) {
			return formatDate(date);
		},
		getCurrentUID() {
			return getCurrentUID();
		},
	},
}
</script>
