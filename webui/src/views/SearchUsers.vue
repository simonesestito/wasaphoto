<script>
export default {
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
    <div>
        <div
            class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 class="h2">Search Users</h1>

            <!-- Username search bar -->
            <div class="input-group w-50">
                <div class="input-group-prepend">
                    <span class="input-group-text" id="search-username-field">@</span>
                </div>
                <input type="text" class="form-control" placeholder="Username" aria-label="Username"
                    aria-describedby="search-username-field">
                <div class="input-group-append">
                    <button class="btn btn-outline-primary" type="button">Search</button>
                </div>
            </div>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
</template>

<style>

</style>
