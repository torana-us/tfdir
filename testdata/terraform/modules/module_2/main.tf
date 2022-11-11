output "name" {
  value = module.child.name
}

module "child" {
  source = "./modules/module_2_1"

}
