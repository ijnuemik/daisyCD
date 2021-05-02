<template>
    <div class="settings">
      <el-container class="container">
        <el-header class="git-header">
            <el-button type="text" @click="httpsFormVisible = true">Connect Repository using HTTPS</el-button>
            <el-dialog title="Connect Repo using HTTPS" :visible.sync="httpsFormVisible">
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="120px">
                <el-form-item label="Repository URL" :label-width="formLabelWidth" prop="url">
                  <el-input placeholder="owner / repository" v-model="ruleForm.url">
                    <template slot="prepend">https://github.com/</template>
                    <template slot="append">.git</template>
                  </el-input>
                </el-form-item>
                <el-form-item label="Username" :label-width="formLabelWidth">
                  <el-input placeholder="" v-model="ruleForm.username">
                  </el-input>
                </el-form-item>
                <el-form-item label="Access token" :label-width="formLabelWidth">
                  <el-input placeholder="required only private repository" v-model="ruleForm.token">
                  </el-input>
                </el-form-item>                    
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="httpsFormVisible = false">Cancel</el-button>
                <el-button type="primary" :loading="loading" v-on:click="confirm()">Confirm</el-button>
            </span>
            </el-dialog>
            <el-button class="ssh-button" type="text" @click="sshFormVisible = true">Connect Repository using SSH</el-button>
            <el-dialog title="Connect Repo using SSH" :visible.sync="sshFormVisible">
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="120px">
                <el-form-item label="Repository URL" :label-width="formLabelWidth" prop="url">
                  <el-input placeholder="Please input repository url" v-model="ruleForm.url">
                    <template slot="prepend">git@github.com:</template>
                    <template slot="append">.git</template>
                  </el-input>
                </el-form-item>
                <el-form-item label="SSH Private Key" :label-width="formLabelWidth">
                  <el-input type="textarea" :rows="5" placeholder="" v-model="ruleForm.sshkey">
                  </el-input>  
                </el-form-item>                 
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="sshFormVisible = false">Cancel</el-button>
                <el-button type="primary" :loading="loading" v-on:click="confirm()">Confirm</el-button>
            </span>
            </el-dialog>
            <el-button type="primary" v-on:click="refresh()"><i class="el-icon-refresh-right el-icon-left">Refresh</i></el-button>
          </el-header>
          <el-main>
            <el-row :gutter="24" v-for="repo in repoList" class="git-row">
              <el-col :span="18" :offset="3" class="git-content">
                <el-col class="git-prop" :span="2">
                  <img class="git-detail-icon" :src="require('@/assets/GitHub-Mark-32px.png')" :size="25"></img>
                </el-col>
                <el-col class="git-prop" :span="2">
                  <el-link class="git-detail-link" :underline="false"> git </el-link>
                </el-col>                
                <el-col class="git-prop" :span="9">
                  <el-link class="git-detail-link" :underline="false"> {{ repo.url }} </el-link>
                </el-col>
                <el-col class="git-prop" :span="4">
                  <el-link class="git-detail-link" :underline="false"> {{ repo.accessuser }} </el-link>
                </el-col> 
                <el-col class="git-prop" :span="5">
                  <el-button class="git-detail-check" type="success" icon="el-icon-check" size="mini" circle></el-button>
                  <el-link class="git-detail-link" :underline="false"> connected </el-link>
                </el-col>
                <el-col class="git-prop" :span="2">
                  <el-popconfirm
                    title="Are you sure to delete this?"
                    confirm-button-text='Yes'
                    cancel-button-text='No'
                    icon="el-icon-info"
                    icon-color="red"
                    @confirm="repoDelete(repo.url)"
                  >
                    <el-button slot="reference" class="git-detail-delete" size="medium" icon="el-icon-delete" circle></el-button>
                  </el-popconfirm>
                </el-col>                                                
              </el-col>
            </el-row>
          </el-main>
      </el-container>
    </div>
</template>

<style>
  .settings{
    flex: 1;
    display: flex;    
    justify-content: center;
    align-items: center;
    margin-top: 0px;
  }
  .git-header{
    display: flex;
    justify-content: center;
    align-items: center;   
    margin-right: auto; 
    margin-left: 0;
  }
  .ssh-button{
    margin-left: 20px;
    margin-right: 20px;
  }
  .git-row{
    margin-top: 20px;
  }
  .git-content{
    background: #d3dce6;
    border-radius: 4px;
    min-height: 60px;  
    display: flex;
    justify-content: center;
    align-items: center;    
  }
  .git-prop{    
    display: flex;
    justify-content: center;
    align-items: center;      
  }
  .git-detail-check{
    margin-right: 10px;
  }
  .git-detail-icon{
    height: 35%;
    width: 35%;
  }
  .git-detail-link{
    font-family: Arial
  }
</style>

<script>
  import axios from 'axios';
  axios.defaults.baseURL = 'http://127.0.0.1:8080';
  const icon = require('@/assets/git-icon.png')
  export default {
    data() {
      return {
        ruleForm: {
          url: '',
          username: '',
          token: '',
          sshkey: ''
        },
        rules: {
          url: [
            { required: true, message: 'Please input repsotiroy URL', trigger: 'blur' }
          ]
        },
        httpsFormVisible: false,
        sshFormVisible: false,
        formLabelWidth: '120px',
        loading: false,
        repoList: [],
        config: {
          headers: {
          'Authorization': 'Bearer ' + this.$cookies.get('accesstoken')
          }
        },
        delete: ''
      };
    },

    created() {
      console.log(this.$cookies.get('accesstoken'))
      axios.get('/repository/list', this.config)
        .then(res => {
          if(res.status == 200){
            console.log(res)
            this.repoList = res.data.data
          }
        })
        .catch(error => {
            console.log('error access')
            return
        })      
    },

    methods: {
      loadingTime() {
        return new Promise(resolve => {
          setTimeout(resolve, 800);
        });
      },
      refresh() {
        console.log(this.$cookies.get('accesstoken'))
        axios.get('/repository/list', this.config)
          .then(res => {
            if(res.status == 200){
              console.log(res)
              this.repoList = res.data.data
            }
          })
          .catch(error => {
              console.log('error access')
              return
          })           
      },
      async confirm(){
        if (this.ruleForm.url == '') {
          return;
        }
        var status = 1
        this.loading = true
        const urlParse = this.ruleForm.url.split('/') // url: github.com/{owner}/{repoName}
        var owner = urlParse[0]
        var repoName = urlParse[1]
        var gitAPI = 'https://api.github.com/repos/'+owner+'/'+repoName
        var repoUrl = 'https://github.com/' +this.ruleForm.url
        await axios.get(gitAPI)
          .then(res => {
            if(res.status == 200) {
              console.log(res)
            }
          }) 
          .catch(error => {
              status = 0 
              this.$message.error('Oops, this is a error message.');
              this.loading = false
              return
          })
        if(status == 1){
          await axios.post('/repository/set', {'url': repoUrl, 'accessuser': this.ruleForm.username, 'accesstoken': this.ruleForm.token, 'status': 1}, this.config)
            .then(res => {
              console.log(res)
            })
            .catch(error => {
              console.log(error)
            })
          this.loading = false          
          this.httpsFormVisible = false
          this.sshFormVisible = false
          this.$message({
            message: 'Congrats, this is a success message.',
            type: 'success'
          });
        }
        this.refresh()
     },
     async repoDelete(url){
       await axios.post('/repository/delete', {'url': url}, this.config)
        .then(res => {
          console.log(res)
        })
        .catch(error => {
          console.log(error)
        })
        this.refresh()
     } 
    },

    computed: {
    }
  };
</script>