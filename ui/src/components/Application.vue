<template>
    <div class="settings">
        <el-container class="container">
          <el-header class="application-header">
            <el-button type="primary" @click="addApplicationFormVisible = true">Add Helm Chart</el-button>
            <el-dialog title="Add Helm Chart" :visible.sync="addApplicationFormVisible">
            <el-form :model="ruleForm" :rules="rules" ref="ruleForm" label-width="120px">
                <el-form-item label="ReleaseName" :label-width="formLabelWidth">
                  <el-input placeholder="must start and end with an alphanumeric character (only lower)" v-model="ruleForm.ReleaseName"></el-input>
                </el-form-item>
                <el-form-item label="Namespace" :label-width="formLabelWidth">
                  <el-input placeholder="" v-model="ruleForm.ReleaseNamespace"></el-input>
                </el-form-item>
                <el-form-item label="Path" :label-width="formLabelWidth">
                  <el-select class="select" placeholder="Select" v-model="ruleForm.HelmChartPath">
                    <el-option v-for="item in helmChartLists" :key="item.value" :label="item.label" :value="item.value">
                    </el-option>
                  </el-select>
                </el-form-item>                    
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addApplicationFormVisible = false">Cancel</el-button>
                <el-button type="primary" :loading="loading" v-on:click="confirm()">Confirm</el-button>
            </span>
            </el-dialog>            
          </el-header>
          <el-main>
            <el-row :gutter="24" v-for="app in appList" class="application-row">
              <el-card class="box-card">
                <div slot="header" class="clearfix">
                  <span>{{ app.releaseName }}</span>
                  <el-button style="float: right; padding: 3px 0" type="text" v-on:click="appDelete(app.releaseName, app.releaseNamespace)">Delete</el-button>
                </div>
                <div class="text item">
                  {{'Release Name: ' + app.releaseName}}
                </div>
                <div class="text item">
                  {{'Release Namespace: ' + app.releaseNamespace}}
                </div>
                <div class="text item">
                  {{'Chart Path: ' + app.chartPath}}
                </div>                                
              </el-card>
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
  .application-header{
    display: flex;
    justify-content: center;
    align-items: center;   
    margin-right: auto; 
    margin-left: 0;
  }
  .select{
    width: 650px;
  }

  .box-card {
    width: 600px;
    margin-bottom: 20px;
    margin-left: 15px;
  } 

  .text {
    font-size: 14px;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
    content: "";
  }
  .clearfix:after {
    clear: both
  }

  .item {
    margin-bottom: 18px;
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
          ReleaseName: '',
          ReleaseNamespace: '',
          HelmChartPath: ''
        },
        rules: {
          ReleaseName: [
            { required: true, message: 'Please input repsotiroy URL', trigger: 'blur' }
          ],
          ReleaseNamespace: [
            { required: true, message: 'Please input repsotiroy URL', trigger: 'blur' }
          ],
          HelmChartPath: [
            { required: true, message: 'Please input repsotiroy URL', trigger: 'blur' }
          ]
        },
        addApplicationFormVisible: false,
        sshFormVisible: false,
        formLabelWidth: '120px',
        loading: false,
        repoList: [],
        config: {
          headers: {
          'Authorization': 'Bearer ' + this.$cookies.get('accesstoken')
          }
        },
        delete: '',
        helmChartLists: [],
        value: '',
        appList: []
      };
    },

    created() {
      console.log(this.$cookies.get('accesstoken'))
      axios.get('/helm/applist', this.config)
        .then(res => {
          if(res.status == 200){
            console.log(res)
            this.appList = res.data.data
          }
        })
        .catch(error => {
          console.log('error')
          return
        })
      axios.get('/helm/chartlist', this.config)
        .then(res => {
          if(res.status == 200){
            console.log(res)
            for(const path of res.data.data){
              var chart = {}
              chart.value = path
              chart.label = path
              console.log(chart.value)
              this.helmChartLists.push(chart)
            }
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
      axios.get('/helm/applist', this.config)
        .then(res => {
          if(res.status == 200){
            console.log(res)
            this.appList = res.data.data
          }
        })
        .catch(error => {
          console.log('error')
          return
        })            
      },

      async confirm(){
        var status = 1
        this.loading = true

        await axios.post('/helm/installchart', {'releaseName': this.ruleForm.ReleaseName, 'releaseNamespace': this.ruleForm.ReleaseNamespace, 'chartPath': this.ruleForm.HelmChartPath}, this.config)
          .then(res => {
            console.log(res)
          })
          .catch(error => {
            console.log(error)
            status = 0
            this.$message.error('Oops, this is a error message.');
            this.loading = false
            return
          })
        if(status == 1){
          this.loading = false          
          this.addApplicationFormVisible = false
          this.sshFormVisible = false
          this.$message({
            message: 'Congrats, helm chart is successfully installed.',
            type: 'success'
          });
        }
        this.refresh()
      },

      async appDelete(appName, appNamespace){
        console.log(appName)

        await axios.post('/helm/appdelete', {'releaseName': appName, 'releaseNamespace': appNamespace}, this.config)
          .then(res => {
            console.log(res)
          })
          .catch(error => {
            console.log(error)
            return
          })  
        this.refresh()  
      }
  },

  computed: {
    }
  };
</script>