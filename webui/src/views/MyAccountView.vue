<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import {saveAuthToken} from "../services/auth-store";
import router from "../router";
export default {
	components: {PageSkeleton},
	data: function () {
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
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</PageSkeleton>
</template>

<style>

</style>
