/**
 * Created by zhangcx on 2017/7/31.
 */

var app = new Vue({
    el: '#app',
    data: {
        key: '',
        value: ''
    },
    methods: {
        queryValue: function () {
            var $this = this
            if (this.key === "") {
                this.value = "The key can not be null!"
                return
            }
            axios.get("/queryValue?key=" + this.key)
                .then(function (response) {
                    $this.key = response.data.Key
                    $this.value = response.data.Value
                }).catch(function (err) {
                $this.value = "Error, could not get the value!" + err
            })
        }
    }
})