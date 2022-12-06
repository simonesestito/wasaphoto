<script>
import PageSkeleton from "../components/PageSkeleton.vue";

export default {
	components: {PageSkeleton},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<PageSkeleton :main-action="{text:'Upload Photo'}" title="Home Page">
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
	</PageSkeleton>
</template>

<style>
</style>
