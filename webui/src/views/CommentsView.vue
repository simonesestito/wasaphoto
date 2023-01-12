<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ShowMore from "../components/ShowMore.vue";
import {CommentsService} from "../services";
import NewCommentModal from "../components/NewCommentModal.vue";
import CommentListItem from "../components/CommentListItem.vue";

export default {
	name: "CommentsView",
	components: {CommentListItem, NewCommentModal, ShowMore, LoadingSpinner, ErrorMsg, PageSkeleton},
	data() {
		return {
			errorMessage: null,
			loading: false,
			pageCursor: null,
			comments: [],
		};
	},
	methods: {
		async refresh() {
			this.pageCursor = null;
			this.comments = [];
			await this.loadMore();
		},
		async loadMore() {
			// Handle double loading or null (/undefined) parameters.
			const photoId = this.$route.params.photoId;
			if (!photoId || this.loading) {
				return;
			}

			this.loading = true;
			this.errorMessage = null;

			try {
				const commentsPage = await CommentsService.getPhotoComments(photoId, this.pageCursor);
				this.pageCursor = commentsPage.nextPageCursor;
				this.comments.push(...commentsPage.pageData);
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async postComment(text) {
			this.loading = true;
			this.errorMessage = null;

			try {
				const newComment = await CommentsService.commentPhoto(
					this.$route.params.photoId,
					text,
				);

				// Add comment to data, no need to refresh the list
				this.comments.unshift(newComment);
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
	},
	created() {
		this.$watch(
			() => this.$route.params,
			(params, _) => {
				this.refresh();
			},
		);
	},
	mounted() {
		this.refresh();
	},
}
</script>

<template>
	<PageSkeleton title="Comments"
				  :main-action="{text: 'Publish comment', onClick: () => $refs.openNewComment.click(),}">
		<ErrorMsg :msg="errorMessage"/>
		<LoadingSpinner v-if="loading"/>

		<!-- Comments list -->
		<CommentListItem v-for="comment in comments" :key="comment.id" :comment="comment" @error="err => errorMessage = err"
						 @refresh="refresh"/>
		<ShowMore v-if="pageCursor" @loadMore="loadMore"/>

		<!-- Empty view -->
		<p v-if="!loading && comments.length === 0 && !errorMessage">
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#frown"/>
			</svg>
			No comments yet here, you can be the first one!
		</p>

		<!-- New comment modal -->
		<NewCommentModal id="new-comment" @comment="postComment"/>
		<!-- Button trigger modal, clicked programmatically -->
		<div ref="openNewComment" data-bs-toggle="modal" data-bs-target="#new-comment"/>
	</PageSkeleton>
</template>
