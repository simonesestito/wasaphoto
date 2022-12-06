<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import {getCurrentUID, saveAuthToken} from "../services/auth-store";
import router from "../router";
import {UsersService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
	components: {LoadingSpinner, PageSkeleton},
	data: function () {
        return {
            errorMessage: null,
            loading: false,
            myProfile: null,
        }
    },
    methods: {
        async refresh() {
            this.loading = true;
            this.errorMessage = null;
            try {
                this.myProfile = await UsersService.getUserProfile(getCurrentUID());
            } catch (e) {
                this.errorMessage = e.toString();
            } finally {
				this.loading = false;
			}
        },
		async onClick() {
			saveAuthToken(null);
			await router.replace('/login');
		}
    },
    mounted() {
        this.refresh()
    }
}
</script>

<template>
	<PageSkeleton title="My Account" :main-action="{text:'Logout', onClick: this.onClick}" :actions="[{text:'Edit account'}]">
		<ErrorMsg v-if="errorMessage" :msg="errorMessage" />

		<LoadingSpinner v-if="loading" />

		<!-- User simple summary -->
		<p v-if="myProfile">Hi, {{ myProfile.name }} {{ myProfile.surname }} (@{{myProfile.username}})</p>
	</PageSkeleton>
</template>

<style>

</style>
