/**
 * Created by zhangcx on 2017/7/31.
 */

var app = new Vue({
    el: '#app',
    data: {
        key: '',
        value: '',
        keyItems: []
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
                    var value = response.data.Value
                    $this.value = isJSON(value) ? JSON.stringify(JSON.parse(value), null, 4) : value
                    if (!$this.keyItems.includes($this.key)) {
                        $this.keyItems.push($this.key)
                    }
                }).catch(function (err) {
                $this.value = "Error, could not get the value!" + err
            })
        },
        saveValue: function () {
            var $this = this
            if (this.key === "" || this.value === ""){
                alert("data is not complete!")
                return
            }
            axios.post("/saveValue", {
                key: $this.key,
                value: $this.value
            }).then(function (response) {
                if (response.data.Msg === "ok"){
                    $this.key = response.data.Key
                    var value = response.data.Value
                    $this.value = isJSON(value) ? JSON.stringify(JSON.parse(value), null, 4) : value
                    return
                }
                alert("save failed!")
            }).catch(function (err) {
                $this.value = "Error, save value to redis failed" + err
            })
        },
        deleteValue: function () {
            var $this = this
            if (this.key === ""){
                alert("key is nil!")
                return
            }
            axios.post("/deleteValue", {
                key: $this.key,
                value: $this.value
            }).then(function (response) {
                if (response.data.Msg === "ok"){
                    $this.key = response.data.Key
                    var value = response.data.Value
                    $this.value = isJSON(value) ? JSON.stringify(JSON.parse(value), null, 4) : value
                    return
                }
            }).catch(function (err) {
                $this.value = "delete value failed: " + err
            })
        }
    }
})

function isJSON(str) {
    if (typeof str == 'string') {
        try {
            JSON.parse(str);
            return true;
        } catch(e) {
            return false;
        }
    }
    console.log('It is not a string!')
}