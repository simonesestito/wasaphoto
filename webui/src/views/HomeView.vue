<script>
import {StreamService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";
import PageSkeleton from "../components/PageSkeleton.vue";
import ShowMore from "../components/ShowMore.vue";
import PhotoListItem from "../components/PhotoListItem.vue";
import router from "../router";

export default {
	name: 'HomeView',
	components: {PhotoListItem, ShowMore, PageSkeleton, ErrorMsg, LoadingSpinner},
	data() {
		return {
			loading: false,
			errorMessage: null,
			photos: [],
			photosCursor: null,
		};
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errorMessage = null;
			this.photos = [];
			this.photosCursor = null;
			await this.loadMore();
		},
		async loadMore() {
			this.loading = true;
			this.errorMessage = null;

			try {
				const photoResponse = await StreamService.getMyStream(this.photosCursor);
				this.photosCursor = photoResponse.nextPageCursor;
				this.photos.push(...photoResponse.pageData);
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async goToUpload() {
			await router.push('/upload');
		},
	},
	mounted() {
		this.refresh();
	},
}
</script>

<template>
	<PageSkeleton title="Home Page" :main-action="{text: 'Upload Photo', onClick: this.goToUpload}">
		<LoadingSpinner v-if="loading"/>
		<ErrorMsg :msg="errorMessage"/>

		<!-- Photos list -->
		<PhotoListItem v-for="photo in photos"
					   :key="photo.id"
					   :photo="photo"
					   @error="(err) => this.errorMessage = err"
					   :show-author="true"/>

		<!-- Empty view -->
		<p v-if="!loading && photos.length === 0 && !errorMessage">
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#frown"/>
			</svg>
			No photos to show, you can
			<RouterLink to="/search">follow someone</RouterLink>
			to see their photos here!
		</p>

		<ShowMore v-if="photosCursor && !loading" @loadMore="loadMore"/>
	</PageSkeleton>
</template>
