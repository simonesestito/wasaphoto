<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import UsernameInput from "../components/UsernameInput.vue";
export default {
	components: {UsernameInput, PageSkeleton},
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
    },
    mounted() {
        this.refresh()
    }
}
</script>

<template>
	<PageSkeleton title="Search Users">
		<!-- Username search bar -->
		<UsernameInput submit-text="Search" />
	</PageSkeleton>
</template>
