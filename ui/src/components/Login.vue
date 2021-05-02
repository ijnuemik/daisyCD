<template>
  <div class="login">
    <el-card>
      <h2>Login</h2>
      <el-form class="login-form" :model="model" :rules="rules" ref="form" @submit.native.prevent="login">
        <el-form-item prop="username">
          <el-input v-model="model.username" placeholder="Username" prefix-icon="el-icon-user"></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="model.password"
            placeholder="Password"
            type="password"
            prefix-icon="el-icon-lock"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" class="login-button" type="primary" native-type="submit" block> Login </el-button>
        </el-form-item>
        <a class="forgot-password" href="https://google.com/">Forgot password ?</a>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: "login",
  data() {
    return {
      model: {
        username: "",
        password: ""
      },
      loading: false,
      rules: {
        username: [
          {
            required: false,
            message: "Username is required",
            trigger: "blur"
          },
          {
            min: 4,
            message: "Username length should be at least 5 characters",
            trigger: "blur"
          }
        ],
        password: [
          { required: false,
          message: "Password is required",
          trigger: "blur" },
          {
            min: 5,
            message: "Password length should be at least 5 characters",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    simulateLogin() {
      return new Promise(resolve => {
        setTimeout(resolve, 800);
      });
    },
    async login() {
      let valid = await this.$refs.form.validate();
      if (!valid) {
        return;
      }
      this.loading = true;
      await this.simulateLogin();
      this.loading = false;

      const formData = {
        username: this.model.username,
        password: this.model.password
      }
      axios.post('http://127.0.0.1:8080/login', formData)
        .then(res => {
          if(res.status == 200) {
            this.$cookies.set("accesstoken", res.data.access_token, 0);
            this.$cookies.set("refreshtoken", res.data.refresh_token, 0);
          }
          this.$message.success("Login successfull");
          console.log(res)
        })
        .catch(error => {
          if (error.response.status == 401) {
            this.$message.error("Username or password is invalid");
          }
        })
    }
  }
};
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .login {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 0px;
    background: rgb(248, 229, 154);
  }
  .login-form {
    margin-top: 50px;
    width: 290px;
    height: 280px;
  }
  .forgot-password {
    margin-top: 10px;
    font-size: x-small;
    color: rgb(186, 187, 177);
  }
  .login-button {
    width: 100%;
    margin-top: 40px;
    background: rgb(241, 196, 0);
    border: rgb(241, 196, 0);
  }
</style>